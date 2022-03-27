package function

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

var testSet = "Allure"
var testLeague = "NHL"

func TestFirestoreIntegration(t *testing.T) {
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	assert.NoError(t, err)
	defer firestoreClient.Close()

	data, err := parseChecklist("test_data.xlsx")

	if assert.NoError(t, err) {
		err := writeToFirestore(testLeague, testSet, data)

		assert.NoError(t, err)

		// Check that there are collections
		setDocument := firestoreClient.Collection(testLeague).Doc(testSet)
		subSets := setDocument.Collections(ctx)

		subCollection, err := subSets.GetAll()

		assert.NoError(t, err)
		assert.NotEmpty(t, subCollection)

		// Check that there are documents
		documents := firestoreClient.Doc(testLeague + "/" + testSet + "/Base Set/checklist")
		documentsSnap, err := documents.Get(ctx)

		assert.NoError(t, err)

		actual := documentsSnap.Data()
		assert.NotEmpty(t, actual)
	}

}
