package ws

import "sync"

type room struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[int64]*room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[int64]*room),
	}
}

func (h *Hub) getOrCreateRoom(roomID int64) *room {
	h.mu.RLock()
	r := h.rooms[roomID]
	h.mu.RUnlock()
	if r != nil {
		return r
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if r = h.rooms[roomID]; r == nil {
		r = &room{clients: make(map[*Client]struct{})}
		h.rooms[roomID] = r
	}
	return r
}

func (h *Hub) JoinRoom(c *Client, roomID int64) {
	r := h.getOrCreateRoom(roomID)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (h *Hub) LeaveRoom(c *Client, roomID int64) {
	h.mu.RLock()
	r := h.rooms[roomID]
	h.mu.RUnlock()
	if r == nil {
		return
	}

	r.mu.Lock()
	delete(r.clients, c)
	empty := len(r.clients) == 0
	r.mu.Unlock()

	if empty {
		h.mu.Lock()
		r.mu.RLock()
		stillEmpty := len(r.clients) == 0
		r.mu.RUnlock()
		if stillEmpty {
			delete(h.rooms, roomID)
		}
		h.mu.Unlock()
	}
}

func (h *Hub) BroadcastToRoom(roomID int64, data []byte, except *Client) {
	h.mu.RLock()
	r := h.rooms[roomID]
	h.mu.RUnlock()
	if r == nil {
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for c := range r.clients {
		if c == except {
			continue
		}
		select {
		case c.send <- data:
		default:
		}
	}
}
