package ApiServer

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/Model/Agent"
	"AITranslatio/app/Model/User"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (s *ApiServer) AddSession(userID int64) (int64, error) {

	//通过UserID找到用户行
	user, err := s.DAO.FindUserByID(userID, "UserID")
	if err != nil {
		return 0, fmt.Errorf("获取用户信息失败！%w", err)
	}

	//序列化其中的Json格式的SessionList
	var sessionList []User.SessionInfo
	if len(user.SessionList) > 0 {
		err = json.Unmarshal(user.SessionList, &sessionList)
		if err != nil {
			return 0, fmt.Errorf("JSON序列化失败%w", err)
		}
	}

	//生成新的雪花ID、新对话 ， 追加到list
	SessionID := s.snowFlakeGenerator.GetID()
	SessionItem := User.SessionInfo{
		Name: "未命名对话",
		ID:   SessionID,
	}
	sessionList = append(sessionList, SessionItem)
	data, err := json.Marshal(sessionList)
	if err != nil {
		return 0, fmt.Errorf("JSON序列化失败%w", err)
	}

	//回写DB
	err = s.DAO.UpdateSessionList(userID, data)
	if err != nil {
		return 0, fmt.Errorf("更新sessionList失败%w", err)
	}

	//同步到redis

	return SessionID, nil
}

// 根据contextID在redis获取上下文
func (s *ApiServer) GetSessionContext(UserID, ContextID string) ([]Agent.ContextItem, error) {

	//使用userID+ContextID拼接，在redis中获取到该context

	key := Consts.RAGKeyPrefix + UserID + ":" + ContextID
	context, err := s.redis.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messagesList []Agent.ContextItem

	for _, v := range context {
		var m Agent.ContextItem
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, err
		}
		messagesList = append(messagesList, m)
	}
	return messagesList, nil
}

func (s *ApiServer) Ask(Question, ContextID, UserID string) (string, error) {

	answer, err := s.RAGClient.Ask(context.Background(), Question, ContextID)
	if err != nil {
		return "", fmt.Errorf("调用RAG服务失败%w", err)
	}

	//更新redis
	key := Consts.RAGKeyPrefix + UserID + ":" + ContextID
	Ask := Agent.ContextItem{
		s.snowFlakeGenerator.GetID(),
		"right",
		"filled",
		Question,
		"MateChat",
		time.Now().String(),
	}
	Answer := Agent.ContextItem{
		s.snowFlakeGenerator.GetID(),
		"left",
		"filled",
		answer,
		"human",
		time.Now().String(),
	}

	//更新redis的contextID对应的内容
	s.redis.RPush(context.Background(), key, Ask, Answer)

	return answer, nil
}
