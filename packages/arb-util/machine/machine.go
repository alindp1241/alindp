package machine

import (
	"github.com/offchainlabs/arbitrum/packages/arb-util/protocol"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
)

type Context interface {
	Send(message value.Value)
	GetTimeBounds() value.Value
	NotifyStep(uint64, bool)
	LoggedValue(value.Value)

	OutMessageCount() int
}

type NoContext struct{}

func (m *NoContext) LoggedValue(data value.Value) {

}

func (m *NoContext) Send(message value.Value) {

}

func (m *NoContext) OutMessageCount() int {
	return 0
}

func (m *NoContext) GetTimeBounds() value.Value {
	return value.NewEmptyTuple()
}

func (m *NoContext) NotifyStep(_ uint64, _ bool) {
}

type Status int

const (
	Extensive Status = iota
	ErrorStop
	Halt
)

type Machine interface {
	Hash() [32]byte
	Clone() Machine
	PrintState()

	CurrentStatus() Status
	LastBlockReason() BlockReason
	InboxHash() value.HashOnlyValue
	DeliverMessages(messages value.TupleValue)

	ExecuteAssertion(maxSteps int32, timeBounds *protocol.TimeBoundsBlocks) (*protocol.ExecutionAssertion, uint32)
	MarshalForProof() ([]byte, error)

	Checkpoint(storage CheckpointStorage) bool
	RestoreCheckpoint(storage CheckpointStorage, checkpointName string) bool
}

func IsMachineBlocked(machine Machine, currentTime uint64) bool {
	lastReason := machine.LastBlockReason()
	if lastReason == nil {
		return false
	}
	return lastReason.IsBlocked(machine, currentTime)
}
