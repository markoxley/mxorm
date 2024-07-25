package mxorm

// Updater is an interface that defines the Update method
//
// The Update method is used to update the model
type Updater interface {
	// Update updates the model
	//
	// This method is called by the Save method of the mxorm.Model
	// interface if the model implements the Updater interface.
	//
	// The implementation of this method should update the model
	// in the database. It is up to the implementation to decide
	// how this is done.
	//
	// For example, the implementation could set the LastUpdate
	// field of the model to the current time, or it could update
	// other fields of the model with new values.
	//
	// This method should return without error if the update is
	// successful, or with an error if the update fails.
	//
	// Implementations of this method should not call the Save
	// method of the mxorm.Model interface.
	Update()
}
