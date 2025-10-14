package WebAuthn

type Vailder struct{}

func (v *Vailder) None() bool {
	return true
}

func (v *Vailder) packed() bool {
	return true
}
