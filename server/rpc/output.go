package rpc

import (
    "fmt"
    "github.com/bytelang/kplayer/module"
    outputype "github.com/bytelang/kplayer/module/output/types"
    "github.com/bytelang/kplayer/proto/msg"
    "github.com/golang/protobuf/proto"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
    "net/http"

    "github.com/bytelang/kplayer/core"
    kpproto "github.com/bytelang/kplayer/proto"
    prompt "github.com/bytelang/kplayer/proto/prompt"
    svrproto "github.com/bytelang/kplayer/server/proto"
)

// Output rpc
type Output struct {
    mm module.ModuleManager
}

func NewOutput(manager module.ModuleManager) *Output {
    return &Output{mm: manager}
}

// Add add output to core player
func (o *Output) Add(r *http.Request, args *svrproto.AddOutputArgs, reply *svrproto.AddOutputReply) error {
    coreKplayer := core.GetLibKplayerInstance()
    if err := coreKplayer.SendPrompt(kpproto.EVENT_PROMPT_ACTION_OUTPUT_ADD, &prompt.EventPromptOutputAdd{
        Path:   []byte(args.Output.Path),
        Unique: []byte(args.Output.Unique),
    }); err != nil {
        return err
    }

    outputModule := o.mm[outputype.ModuleName]
    outputAddMsg := &msg.EventMessageOutputAdd{}

    keeperCtx := module.NewKeeperContext(uuid.New().String(), kpproto.EVENT_MESSAGE_ACTION_OUTPUT_ADD, func(msg []byte) bool {
        if err := proto.Unmarshal(msg, outputAddMsg); err != nil {
            log.Fatal(err)
        }
        return string(outputAddMsg.Unique) == args.Output.Unique
    })
    defer keeperCtx.Close()

    if err := outputModule.RegisterKeeperChannel(keeperCtx); err != nil {
        return err
    }

    // wait context
    keeperCtx.Wait()

    if outputAddMsg.Error != nil {
        log.Errorf("%s", outputAddMsg.Error)
        return fmt.Errorf("%s", outputAddMsg.Error)
    }

    reply.Output.Path = string(outputAddMsg.Path)
    reply.Output.Unique = string(outputAddMsg.Unique)

    return nil
}
