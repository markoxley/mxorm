package mxorm

type Updater interface {
	// Update updates the model
	Update()
}

type Restorer interface {
	Restore()
}
