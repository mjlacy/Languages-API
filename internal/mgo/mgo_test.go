package mgo

//func Test_Find_ShouldReturnLanguages(t *testing.T) {
//	mockCollection := mockCollection{cursor: mongoCursor, err: nil}
//
//	coll := MongoCollection{Collection: mockCollection}
//
//	langs, _ := coll.Find(models.Language{})
//
//	if !reflect.DeepEqual(cursor.Cursor, mongoCursor) {
//		t.Error("Cursor does not match")
//	}
//}
//
//func Test_Find_ShouldReturnError(t *testing.T) {
//	firstAppeared := time.Date(2009, 11, 10, 0, 0, 0, 0, time.UTC)
//	lang := models.Language{
//		Name: "Golang",
//		Creators: []string{
//			"Robert Griesemer",
//			"Rob Pike",
//			"Ken Thompson",
//		},
//		Extensions: []string{
//			".go",
//		},
//		FirstAppeared: &firstAppeared,
//		Year:          2009,
//		Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
//	}
//	err := errors.New("error")
//	mockCollection := mockCollection{cursor: nil, err: err}
//
//	coll := MongoCollection{Collection: mockCollection}
//
//	_, findErr := coll.Find(lang)
//
//	if !errors.Is(err, findErr) {
//		t.Error("Error does not match")
//	}
//}
//
//func Test_FindOne_ShouldReturnMongoSingleResultWithErrorIfGivenInvalidId(t *testing.T) {
//	mockCollection := mockCollection{err: errors.New("error")}
//
//	coll := MongoCollection{Collection: mockCollection}
//
//	msr := coll.FindOne("Invalid")
//
//	result, ok := msr.(MongoSingleResult)
//	if !ok {
//		t.Error("Result is not a MongoSingleResult")
//	}
//
//	if !errors.Is(result.error, models.ErrInvalidId) {
//		t.Error("Errors do not match")
//	}
//}
//
//func Test_FindOne_ShouldReturnMongoSingleResult(t *testing.T) {
//	sr := &mongo.SingleResult{}
//	mockCollection := mockCollection{singleResult: sr}
//
//	coll := MongoCollection{Collection: mockCollection}
//
//	msr := coll.FindOne("6657d0394fb1d5acb908f30f")
//
//	result, ok := msr.(MongoSingleResult)
//	if !ok {
//		t.Error("Result is not a MongoSingleResult")
//	}
//
//	if !reflect.DeepEqual(result.SingleResult, sr) {
//		t.Error("Errors do not match")
//	}
//}
