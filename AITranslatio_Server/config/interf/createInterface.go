package interf

type CreateInterf interface {
	Create(fileName string) ConfigInterface
}

func CreateConfig(interf CreateInterf, fileName string) ConfigInterface {
	return interf.Create(fileName)
}
