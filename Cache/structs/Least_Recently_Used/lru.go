package Least_Recently_Used

import (
	"sync"
	"time"
)

// Пример какой-то сущности
type Order struct {
	id   int
	name string
}

// Связной список
type Node struct {
	key       string    // key — уникальный идентификатор заказа (например, order_uid).
	value     *Order    // value — данные заказа (*Order), именно это мы кешируем
	expiresAt time.Time // expiresAt — момент времени, когда запись "протухнет" (для TTL).
	prev      *Node     // prev — ссылка на предыдущий элемент
	next      *Node     // next — ссылка на следующий элемент.
}

/*
Так легко вставлять и удалять элементы из середины за O(1).

Без списка пришлось бы "сдвигать массив" (например, queue []string), что медленно.
*/

// Структра нашего кеша
type LRUCache struct {
	capacity int              // Емкость кеша
	data     map[string]*Node // map для того чтобы добиться (O(1)) в Get.
	head     *Node            // фиктивная голова
	tail     *Node            // фиктивный хвост Без них пришлось бы постоянно проверять «а вдруг список пустой?», «а вдруг это последний элемент?».
	mu       sync.Mutex       // для защиты от data race.
	ttl      time.Duration
}

// NewLRUCache : Новый кеш.
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	l := &LRUCache{
		capacity: capacity,
		data:     make(map[string]*Node),
		ttl:      ttl,
	}
	l.head = &Node{}
	l.tail = &Node{}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

// Get : возвращает заказ по ключу, если он есть и не протух по TTL
func (l *LRUCache) Get(key string) *Order {
	l.mu.Lock()
	defer l.mu.Unlock() // RWMutex имеет смысл только если у тебя есть методы, которые реально только читают, а у нас Get изменяет порядок → значит нужен Lock.

	node, ok := l.data[key]
	if ok {
		// проверяем TTL
		if time.Now().After(node.expiresAt) {
			// удаляем протухший элемент
			l.removeNode(node)
			delete(l.data, key)
			return nil
		}
		// делаем "свежим"
		l.moveToHead(node)
		return node.value
	}
	return nil
}

// Set кладёт заказ в кеш (обновляет, если уже был)
func (l *LRUCache) Set(key string, value *Order) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, ok := l.data[key]
	if ok {
		// обновляем значение и TTL
		node.value = value
		node.expiresAt = time.Now().Add(l.ttl)
		l.moveToHead(node)
		return
	}

	// если переполнено — удаляем самый старый
	if len(l.data) >= l.capacity {
		old := l.removeTail()
		delete(l.data, old.key)
	}

	// создаём новый узел
	newNode := &Node{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(l.ttl),
	}
	l.data[key] = newNode
	l.addToHead(newNode)
}

// addToHead : Функция вставляет новый узел сразу после HEAD.
func (l *LRUCache) addToHead(n *Node) {
	n.prev = l.head      // связываем новый узел с головой
	n.next = l.head.next // новый указывает на бывший первый элемент
	l.head.next.prev = n // бывший первый элемент теперь знает, что перед ним node
	l.head.next = n      // голова теперь смотрит на новый элемент
}

// removeNode : Удалить узел из списка.
func (l *LRUCache) removeNode(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// removeTail : Удалить наименее недавно использованный элемент.
func (l *LRUCache) removeTail() *Node {
	n := l.tail.prev
	l.removeNode(n)
	return n
}

// moveToHead : «освежить» позицию элемента при использовании (убрать из старого места и перенести в начало).
func (l *LRUCache) moveToHead(n *Node) {
	l.removeNode(n)
	l.addToHead(n)
}
