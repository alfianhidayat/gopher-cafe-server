package worker

import (
	"errors"
	"gopher-cafe/internal/entity/coffeeshop"
	"sync"

	"github.com/ajaibid/coin-common-golang/logger"
)

var (
	EquipmentWorkers = map[coffeeshop.EquipmentType]uint8{
		coffeeshop.EquipEspressoMachine: 2,
		coffeeshop.EquipGrinder:         1,
		coffeeshop.EquipMilkSteamer:     1,
		coffeeshop.EquipBlender:         1,
		coffeeshop.EquipWhisk:           2,
	}
)

type EquipPoolManager struct {
	pools map[coffeeshop.EquipmentType]*WorkerPool
	mu    sync.RWMutex
}

func NewEquipPoolManager(totalPool uint8) *EquipPoolManager {
	return &EquipPoolManager{
		pools: make(map[coffeeshop.EquipmentType]*WorkerPool, totalPool),
	}
}

func (e *EquipPoolManager) Register(equipType coffeeshop.EquipmentType, numOfWorkers uint8) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.pools[equipType] = NewWorkerPool(equipType.String(), numOfWorkers)
}

func (e *EquipPoolManager) GetWorkerPool(equipType coffeeshop.EquipmentType) (*WorkerPool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if pool, ok := e.pools[equipType]; ok {
		return pool, nil
	}

	return nil, errors.New("no worker pool registered")
}

func (e *EquipPoolManager) StartAll() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, pool := range e.pools {
		pool.start()
	}

	logger.Debugf("All workers started")
}

func (e *EquipPoolManager) StopAll() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, pool := range e.pools {
		pool.stop()
	}

	logger.Debugf("All workers stopped")
}
