package rpc

import (
    "github.com/bytelang/kplayer/module"
    resourcetype "github.com/bytelang/kplayer/module/resource/types"
    "github.com/bytelang/kplayer/proto/msg"
    "github.com/golang/protobuf/proto"
    "github.com/google/uuid"
    "net/http"

    "github.com/bytelang/kplayer/core"
    kpproto "github.com/bytelang/kplayer/proto"
    kpprompt "github.com/bytelang/kplayer/proto/prompt"
    svrproto "github.com/bytelang/kplayer/server/proto"
    log "github.com/sirupsen/logrus"
)

// Resource rpc
type Resource struct {
    mm module.ModuleManager
}

func NewResource(manager module.ModuleManager) *Resource {
    return &Resource{mm: manager}
}

// Add add Resource to core
func (s *Resource) Add(r *http.Request, args *svrproto.AddResourceArgs, reply *svrproto.AddResourceReply) error {
    coreKplayer := core.GetLibKplayerInstance()
    if err := coreKplayer.SendPrompt(kpproto.EVENT_PROMPT_ACTION_RESOURCE_ADD, &kpprompt.EventPromptResourceAdd{
        Path:   []byte(args.Res.Path),
        Unique: []byte(args.Res.Unique),
    }); err != nil {
        return err
    }

    resourceModule := s.mm[resourcetype.ModuleName]
    resourceAddMsg := &msg.EventMessageResourceAdd{}

    keeperCtx := module.NewKeeperContext(uuid.New().String(), kpproto.EVENT_MESSAGE_ACTION_RESOURCE_ADD, func(msg []byte) bool {
        if err := proto.Unmarshal(msg, resourceAddMsg); err != nil {
            log.Fatal(err)
        }
        return string(resourceAddMsg.Unique) == args.Res.Unique
    })
    defer keeperCtx.Close()

    if err := resourceModule.RegisterKeeperChannel(keeperCtx); err != nil {
        return err
    }

    // wait context
    keeperCtx.Wait()

    reply.Res.Unique = string(resourceAddMsg.Unique)
    reply.Res.Path = string(resourceAddMsg.Path)

    return nil
}

// Remove remove Resource to core
func (s *Resource) Remove(r *http.Request, args *svrproto.RemoveResourceArgs, reply *svrproto.RemoveResourceReply) error {
    coreKplayer := core.GetLibKplayerInstance()
    if err := coreKplayer.SendPrompt(kpproto.EVENT_PROMPT_ACTION_RESOURCE_REMOVE, &kpprompt.EventPromptResourceRemove{
        Unique: []byte(args.Unique),
    }); err != nil {
        return err
    }

    ResourceModule := s.mm[resourcetype.ModuleName]
    resourceRemoveMsg := &msg.EventMessageResourceRemove{}

    keeperCtx := module.NewKeeperContext(uuid.New().String(), kpproto.EVENT_MESSAGE_ACTION_RESOURCE_REMOVE, func(msg []byte) bool {
        if err := proto.Unmarshal(msg, resourceRemoveMsg); err != nil {
            log.Fatal(err)
        }
        return string(resourceRemoveMsg.Unique) == args.Unique
    })
    defer keeperCtx.Close()

    if err := ResourceModule.RegisterKeeperChannel(keeperCtx); err != nil {
        return err
    }

    // wait context
    keeperCtx.Wait()

    reply.Res.Unique = string(resourceRemoveMsg.Unique)
    reply.Res.Path = string(resourceRemoveMsg.Path)

    return nil
}
