package replication

import (
	"context"

	"concurrency/app/config"
	"concurrency/app/network"
	"go.uber.org/zap"
)

const (
	MasterType = "master"
	SlaveType  = "slave"
)

type Replica struct {
	Master *Master
	Slave  *Slave
}

func NewReplication(ctx context.Context, replicationCfg *config.Replication, walCfg *config.WAL, logger *zap.Logger) (*Replica, error) {
	if replicationCfg == nil || walCfg == nil {
		return nil, nil
	}
	replica := &Replica{}
	replicationType := replicationCfg.ReplicaType

	if replicationType == MasterType {
		cfgNetwork := &config.Network{
			Address:        replicationCfg.MasterAddress,
			MaxMessageSize: walCfg.MaxSegmentSize,
		}

		server, err := network.NewTCPServer(ctx, cfgNetwork, logger)
		if err != nil {
			return nil, err
		}
		master := NewMaster(server, walCfg.DataDirectory, logger)
		replica.Master = master
	}
	if replicationType == SlaveType {
		client, err := network.NewTcpClient(replicationCfg.MasterAddress)
		if err != nil {
			return nil, err
		}
		slave, err := NewSlave(client, walCfg.DataDirectory, replicationCfg.SyncInterval, logger)
		if err != nil {
			return nil, err
		}
		replica.Slave = slave
	}
	return replica, nil
}
