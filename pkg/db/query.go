package db

// Common Data Store Query
func (store CommonDataStore) SelectAll(beanList interface{}) error {
	return store.ormEngine.Find(beanList)
}
