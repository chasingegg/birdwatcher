package states

import (
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/milvus-io/birdwatcher/framework"
	"github.com/milvus-io/birdwatcher/models"
	"github.com/milvus-io/birdwatcher/proto/v2.0/querypb"
	querypbv2 "github.com/milvus-io/birdwatcher/proto/v2.2/querypb"
)

type queryCoordState struct {
	*framework.CmdState
	session   *models.Session
	client    querypb.QueryCoordClient
	clientv2  querypbv2.QueryCoordClient
	conn      *grpc.ClientConn
	prevState framework.State
}

// SetupCommands setups the command.
// also called after each command run to reset flag values.
func (s *queryCoordState) SetupCommands() {
	cmd := &cobra.Command{}
	cmd.AddCommand(
		// metrics
		getMetricsCmd(s.client),
		// configuration
		getConfigurationCmd(s.clientv2, s.session.ServerID),
		// back
		getBackCmd(s, s.prevState),
		// exit
		getExitCmd(s),
	)
	s.MergeFunctionCommands(cmd, s)

	s.CmdState.RootCmd = cmd
	s.SetupFn = s.SetupCommands
}

/*
func (s *queryCoordState) ShowCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "show-collection",
		Run: func(cmd *cobra.Command, args []string) {
			collection, err := cmd.Flags().GetInt64("collection")
			if err != nil {
				cmd.Usage()
				return
			}

			req := &querypbv2.ShowCollectionsRequest{
				Base: &commonpbv2.MsgBase{
					TargetID: s.session.ServerID,
				},
				CollectionIDs: []int64{collection},
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			resp, err := s.clientv2.ShowCollections(ctx, req)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Printf("%s, %s", resp.GetStatus().GetErrorCode().String(), resp.GetStatus().GetReason())
		},
	}

	cmd.Flags().Int64("collection", 0, "collection to show")
	return cmd
}*/

func getQueryCoordState(client querypb.QueryCoordClient, conn *grpc.ClientConn, prev framework.State, session *models.Session) framework.State {
	state := &queryCoordState{
		CmdState:  framework.NewCmdState(fmt.Sprintf("QueryCoord-%d(%s)", session.ServerID, session.Address)),
		session:   session,
		client:    client,
		clientv2:  querypbv2.NewQueryCoordClient(conn),
		conn:      conn,
		prevState: prev,
	}

	state.SetupCommands()

	return state
}
