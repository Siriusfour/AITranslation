package BaseDAO

import (
	"AITranslatio/Src/Model"
)

func (BaseDAO *BaseDAO) CreateNovelProgramming(NoteInfo *Model.Note) error {

	result := BaseDAO.orm.Create(NoteInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (BaseDAO *BaseDAO) CreateBranch(BranchInfo *Model.Branch) error {
	result := BaseDAO.orm.Create(BranchInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (BaseDAO *BaseDAO) CreateCommit(CommitInfo *Model.Commit) error {
	result := BaseDAO.orm.Create(CommitInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (BaseDAO *BaseDAO) FindNote(NoteID int) (Model.Note, error) {

	Note := &Model.Note{}

	result := BaseDAO.orm.First(Note, NoteID)
	if result.Error != nil {
		return Model.Note{}, result.Error
	}

	return *Note, nil
}

func (BaseDAO *BaseDAO) FindBranches(NoteID int) (*[]Model.Branch, error) {

	BranchInfo := &[]Model.Branch{}
	result := BaseDAO.orm.Where("NoteID = ?", NoteID).Find(BranchInfo)
	if result.Error != nil {
		return BranchInfo, result.Error
	}

	return BranchInfo, nil
}

func (BaseDAO *BaseDAO) FindCommit(BranchID uint) (*[]Model.Commit, error) {
	Commits := &[]Model.Commit{}
	result := BaseDAO.orm.Where("BranchID = ?", BranchID).Find(Commits)
	if result.Error != nil {
		return Commits, result.Error
	}
	return Commits, nil
}

// ChangeCommit Todo：修改权限，通过邮箱验证码确认过身份猴
func (BaseDAO *BaseDAO) ChangeCommit(CommitID uint) (*Model.Commit, int, int, error) {
	CommitInfo := &Model.Commit{}
	BranchInfo := &Model.Branch{}
	result := BaseDAO.orm.Where("CommitID = ?", CommitID).Find(CommitInfo)
	if result.Error != nil {
		return &Model.Commit{}, 0, 0, result.Error
	}

	result = BaseDAO.orm.Where("BranchID = ?", CommitInfo.BranchID).Find(BranchInfo)
	if result.Error != nil {
		return &Model.Commit{}, 0, 0, result.Error
	}

	return CommitInfo, BranchInfo.Permissions, BranchInfo.WriterID, nil
}

func (BaseDAO *BaseDAO) ChangeFile(CommitID uint, Filepath string) error {

	result := BaseDAO.orm.Model(Model.Commit{}).Where("CommitID", CommitID).Update("FilePath", Filepath)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
