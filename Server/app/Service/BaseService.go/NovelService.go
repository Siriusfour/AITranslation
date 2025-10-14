package BaseService

import (
	"AITranslatio/app/Model"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/DTO"
	"AITranslatio/http/reposen"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strconv"
)

func (BaseService *AuthService.BaseService) CreateNovelProgramming(NovelDTO *DTO.NovelDTO) error {

	NewNovel := Model.Note{
		WriterID:     NovelDTO.WriterID,     // 作者ID
		BranchCount:  0,                     // 分支数
		Introduction: NovelDTO.Introduction, // 简介
		ReaderCount:  0,                     // 订阅人数
		NoteName:     NovelDTO.NoteName,
		Permissions:  NovelDTO.Permissions,
	}

	err := BaseService.BaseDAO.CreateNovelProgramming(&NewNovel)
	if err != nil {
		return err
	}

	NewBranch := Model.Branch{
		WriterID:     NovelDTO.WriterID,
		ReaderCount:  0,
		Introduction: NovelDTO.Introduction,
		Permissions:  NovelDTO.Permissions,
		BranchName:   "Master",
		NoteID:       NewNovel.ID,
	}

	//创建一条主分支
	err = BaseService.BaseDAO.CreateBranch(&NewBranch)
	if err != nil {
		return err
	}

	return nil

}

func (BaseService *AuthService.BaseService) CreateBranch(BranchDTO *DTO.Branch) error {

	NewBranch := Model.Branch{
		WriterID:     BranchDTO.WriterID,     // 作者ID 		// 分支数
		Introduction: BranchDTO.Introduction, // 简介
		ReaderCount:  0,                      // 订阅人数
		BranchName:   BranchDTO.BranchName,
		Permissions:  BranchDTO.Permissions,
		CommitID:     BranchDTO.CommitID, //起始ID
	}

	err := BaseService.BaseDAO.CreateBranch(&NewBranch)
	if err != nil {
		return err
	}
	return nil

}

func (BaseService *AuthService.BaseService) CreateCommit(commitCTX *gin.Context, CommitDTO *DTO.CommitDTO) error {

	//接受文件保存到upload文件夹
	file, err := commitCTX.FormFile("file")
	if err != nil {
		return err
	}

	filename := strconv.Itoa(CommitDTO.WriterID) + "_" + filepath.Base(file.Filename)
	Filepath := filepath.Join("UpLoad", filename)

	if err := commitCTX.SaveUploadedFile(file, Filepath); err != nil {
		println(err.Error())
		return err
	}

	CommitInfo := Model.Commit{
		WriterID:     CommitDTO.WriterID,
		Introduction: CommitDTO.Introduction,
		CommitName:   CommitDTO.CommitName,
		BranchID:     CommitDTO.BranchID,
		FilePath:     Filepath,
		LastNode:     CommitDTO.LastNode,
		NextNode:     CommitDTO.NextNode,
	}

	err = BaseService.BaseDAO.CreateCommit(&CommitInfo)
	if err != nil {
		return err
	}

	return nil
}

func (BaseService *AuthService.BaseService) Programming(NoteID int) (reposen.Note, error) {

	HTTPNote := reposen.Note{}
	HTTPNote.Branches = []reposen.Branches{}

	//查询该项目的记录
	note, err := BaseService.BaseDAO.FindNote(NoteID)
	if err != nil {
		return reposen.Note{}, err
	}
	HTTPNote.Note = &note

	//查询该项目的所有分支,	//循环查询这些分支的所有提交
	Branches, err := BaseService.BaseDAO.FindBranches(NoteID)
	for _, Branch := range *Branches {
		Commits, err := BaseService.BaseDAO.FindCommit(Branch.ID)
		if err != nil {
			return reposen.Note{}, err
		}
		HTTPNote.Branches = append(HTTPNote.Branches, reposen.Branches{
			Branch:  &Branch,
			Commits: Commits,
		})

	}

	//组合成数据返回
	return HTTPNote, nil
}

func (BaseService *AuthService.BaseService) ChangeCommit(Ctx *gin.Context, CommitDTO *DTO.CommitDTO) error {

	//查询数据
	CommitInfo, Permissions, WriterID, err := BaseService.BaseDAO.ChangeCommit(CommitDTO.CommitID)
	if err != nil {
		return err
	}

	//检查提交者权限
	if Permissions == 2 {
		if WriterID != CommitDTO.WriterID {
			return errors.New("权限不足")
		}
	}

	//保存文件，更改提交指向的文件位置,删除原本的数据
	file, err := Ctx.FormFile("file")
	if err != nil {
		return err
	}

	filename := strconv.Itoa(CommitDTO.WriterID) + "_" + filepath.Base(file.Filename)
	Filepath := filepath.Join("UpLoad", filename)
	err = BaseService.BaseDAO.ChangeFile(CommitDTO.CommitID, Filepath)
	if err != nil {
		return err
	}

	err = os.Remove(CommitInfo.FilePath)
	if err != nil {
		return err
	}

	return nil

}
