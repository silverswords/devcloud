package workflows

import (
	"context"
	"log"

	"github.com/cschleiden/go-workflows/backend/sqlite"
	"github.com/cschleiden/go-workflows/client"
	"github.com/cschleiden/go-workflows/worker"
	"github.com/cschleiden/go-workflows/workflow"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/silverswords/devcloud/pkg/appctx"
	"github.com/silverswords/devcloud/pkg/tasks"
)

type Config struct {
	Enabled   bool
	Name      string
	Task      string
	Frequency string
}

var taskDB = sqlite.NewInMemoryBackend()

func RunWorker() {
	var config []Config
	if data := appctx.GlobalAppContext.Viper.Get("tasks"); data == nil {
		return
	} else {
		mapstructure.Decode(data, &config)
	}
	log.Println(appctx.GlobalAppContext.Viper.Get("tasks"))

	w := worker.New(taskDB, nil)

	w.RegisterWorkflow(MongoWorkflow)
	w.RegisterActivity(tasks.DockerRunMongo)

	if err := w.Start(context.TODO()); err != nil {
		panic("could not start worker")
	}

	startWorkflow(config[0])
}

func startWorkflow(config Config) {
	if config.Enabled {
		c := client.New(taskDB)
		_, err := c.CreateWorkflowInstance(context.Background(), client.WorkflowInstanceOptions{
			InstanceID: uuid.NewString(),
		}, MongoWorkflow, "Hello world")
		if err != nil {
			panic("could not start workflow")
		}

	}
}

func MongoWorkflow(ctx workflow.Context, msg string) (string, error) {
	logger := workflow.Logger(ctx)
	logger.Debug("Entering MongoWorkflow", msg)

	if r0, err := workflow.ExecuteActivity[string](ctx, workflow.DefaultActivityOptions, tasks.DockerRunMongo, "hello").Get(ctx); err != nil {
		logger.Debug("error getting activity success result", err)
	} else {
		logger.Debug("ActivitySuccess result:", r0)
	}

	return "", nil
}
