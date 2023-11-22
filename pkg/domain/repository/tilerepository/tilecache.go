package tilerepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/graphtile"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"sync"
)

const cacheMissesBeforeLogging = 1000
const cacheOverflowsBeforeLogging = 1000

type tileCache interface {
	Get(name string) *graphtile.GraphTile
	Set(name string, tile *graphtile.GraphTile)
	SetCacheSize(maxCacheSize int)
	WriteToDisk()
}

type tileCacheImpl struct {
	logger         logging.Logger
	cacheLock      sync.RWMutex
	cache          *listItem
	reader         tileReader
	writer         tileWriter
	maxCacheSize   int
	cacheMisses    int
	cacheOverflows int
}

type listItem struct {
	name string
	tile *graphtile.GraphTile
	next *listItem
}

func newTileCache(maxCacheSize int, reader tileReader, writer tileWriter, logger logging.Logger) tileCache {
	return &tileCacheImpl{
		logger:         logger,
		reader:         reader,
		writer:         writer,
		cache:          nil,
		maxCacheSize:   maxCacheSize,
		cacheMisses:    0,
		cacheOverflows: 0,
	}
}

func (t *tileCacheImpl) Get(name string) *graphtile.GraphTile {
	item := t.find(name)
	if item != nil {
		return item.tile
	}

	t.cacheMisses++
	if t.cacheMisses%cacheMissesBeforeLogging == 0 {
		t.logger.Info().Msgf("cache misses: %d", t.cacheMisses)
	}

	tile, err := t.reader.Read(name)
	if err != nil {
		return nil
	}

	t.insertFront(name, tile)
	return tile
}

func (t *tileCacheImpl) Set(name string, tile *graphtile.GraphTile) {
	item := t.find(name)
	if item != nil {
		item.tile = tile
		return
	}

	t.insertFront(name, tile)
}

func (t *tileCacheImpl) insertFront(name string, tile *graphtile.GraphTile) {
	item := &listItem{
		name: name,
		tile: tile,
		next: t.cache,
	}

	t.cacheLock.Lock()
	defer t.cacheLock.Unlock()

	t.cache = item

	go t.cacheOverflowHandler()
}

func (t *tileCacheImpl) find(name string) *listItem {
	t.cacheLock.RLock()
	defer t.cacheLock.RUnlock()

	var prev *listItem
	current := t.cache
	for current != nil {
		if current.name == name {
			// reset item to front if it got accessed
			if prev != nil {
				prev.next = current.next
				current.next = t.cache
				t.cache = current
			}
			return current
		}
		prev = current
		current = current.next
	}
	return nil
}

func (t *tileCacheImpl) cacheOverflowHandler() {
	t.cacheLock.RLock()
	var prev *listItem
	current := t.cache
	// skip allowed cache size
	for i := 0; i < t.maxCacheSize && current != nil; i++ {
		prev = current
		current = current.next
	}
	t.cacheLock.RUnlock()

	t.cacheLock.Lock()
	if current != nil {
		prev.next = nil
	}
	defer t.cacheLock.Unlock()

	for current != nil {
		t.cacheOverflows++
		if t.cacheOverflows%cacheOverflowsBeforeLogging == 0 {
			t.logger.Info().Msgf("cache overflows: %d", t.cacheOverflows)
		}

		err := t.writer.Write(current.name, current.tile)
		if err != nil {
			t.logger.Error().Msgf("could not write tile %s to disk: %s", current.name, err.Error())
		}
		current = current.next
	}
}

func (t *tileCacheImpl) WriteToDisk() {
	current := t.cache

	t.cacheLock.Lock()
	t.cache = nil
	defer t.cacheLock.Unlock()

	for current != nil {
		err := t.writer.Write(current.name, current.tile)
		if err != nil {
			t.logger.Error().Msgf("could not write tile %s to disk: %s", current.name, err.Error())
		}
		current = current.next
	}
}

func (t *tileCacheImpl) SetCacheSize(maxCacheSize int) {
	t.cacheLock.Lock()
	defer t.cacheLock.Unlock()

	t.maxCacheSize = maxCacheSize
}
