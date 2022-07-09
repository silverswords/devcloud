package main

import (
	"context"
	"os"
	"time"

	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
	_ "gocloud.dev/docstore/mongodocstore"

	"github.com/sirupsen/logrus"
	"gocloud.dev/docstore"
	"gocloud.dev/gcerrors"
)

type Player struct {
	Name             string `docstore:"name,omitempty"`
	Score            int
	DocstoreRevision interface{}
}

func main() {
	ctx := context.Background()

	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	os.Setenv("MONGO_SERVER_URL", "mongodb://localhost:27017")

	coll, err := docstore.OpenCollection(context.Background(), "mongo://my-db/my-coll?id_field=name")
	if err != nil {
		logrus.Errorf("could not open collection: ", err)
		return
	}
	defer coll.Close()

	pat := &Player{Name: "Pat", Score: 7}
	// Clear player
	if err := coll.Delete(ctx, pat); err != nil {
		logrus.Info(err)
	}

	// Create a player.
	if err := coll.Create(ctx, pat); err != nil {
		logrus.Info(err)
	}
	logrus.Printf("%+v\n", pat) // memdocstore revisions are deterministic, so we can check the output.

	// Double a player's score. We cannot use Update to multiply, so we use optimistic
	// locking instead.

	// We may have to retry a few times; put a time limit on that.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	logrus.Println(pat)

	defer cancel()
	for {
		// Get the document.
		player := &Player{Name: "Pat"}
		if err := coll.Get(ctx, player); err != nil {
			logrus.Info(err)
		}
		// player.DocstoreRevision is set to the document's revision.

		// Modify the document locally.
		player.Score *= 2

		// Replace the document. player.DocstoreRevision will be checked against
		// the stored document's revision.
		err := coll.Replace(ctx, player)
		if err != nil {
			code := gcerrors.Code(err)
			// On FailedPrecondition or NotFound, try again.
			if code == gcerrors.FailedPrecondition || code == gcerrors.NotFound {
				continue
			}
			logrus.Info(err)
		}
		break
	}
}
