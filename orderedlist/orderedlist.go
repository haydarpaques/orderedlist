/*
	Package for storing pairs of key-value where the keys are of type string
	and the values are of type 64-bit unsigned integer. The list of pairs are ordered
	starting from the lowest value down to the highest.
*/
package orderedlist

import (
	"errors"
)

type (
	Record struct {
		Key   string
		Value uint64
	}

	OrderedList struct {
		Bookkeeping map[string]bool
		Rec         []Record
	}
)

// New returns new (empty) instance of ordered list
func New() OrderedList {
	return OrderedList{
		Bookkeeping: make(map[string]bool),
	}
}

// Insert inserts a key-value pair into the struct
func (ol *OrderedList) Insert(key string, value uint64) error {
	exists := ol.Bookkeeping[key]
	if exists {
		return errors.New("key already exists")
	}

	// Get proper index
	index := ol.getPosition(value)

	if index == -1 {
		// Append to the bottom of the list
		ol.Rec = append(ol.Rec, Record{key, value})
	} else {
		// Append accordingly
		ol.Rec = append(ol.Rec, Record{})
		copy(ol.Rec[index+1:], ol.Rec[index:])
		ol.Rec[index] = Record{key, value}
	}

	ol.Bookkeeping[key] = true
	return nil
}

// Get returns value of the provided key
func (ol *OrderedList) Get(key string) (uint64, error) {
	index, err := ol.getIndexByKey(key)
	if err != nil {
		return 0, err
	}

	return ol.Rec[index].Value, nil
}

// GetLowest returns key-value pair with lowest value on the list
func (ol *OrderedList) GetLowest() uint64 {
	return ol.Rec[0].Value
}

// GetHighest returns key-value pair with highest value on the list
func (ol *OrderedList) GetHighest() uint64 {
	return ol.Rec[len(ol.Rec)-1].Value
}

// Remove removes a key-value pair by key
func (ol *OrderedList) Remove(key string) error {
	index, err := ol.getIndexByKey(key)
	if err != nil {
		return err
	}

	ol.Rec = append(ol.Rec[:index], ol.Rec[index+1:]...)
	delete(ol.Bookkeeping, key)

	return nil
}

// Update updates value by key
func (ol *OrderedList) Update(key string, value uint64) error {
	if err := ol.Remove(key); err != nil {
		return err
	}

	if err := ol.Insert(key, value); err != nil {
		return err
	}

	return nil
}

// Merge merges list passed to the parameter into the receiver.
// Source list shall be left unchanged.
// Duplicate key error will not be reported.
func (ol *OrderedList) Merge(source *OrderedList) {
	for _, record := range source.Rec {
		ol.Insert(record.Key, record.Value)
	}
}

// getPosition returns proper position (index) to be placed into struct
func (ol *OrderedList) getPosition(value uint64) int {
	index := -1
	for i, record := range ol.Rec {
		if value <= record.Value {
			index = i
			break
		}
	}

	return index
}

// getIndexByKey returns index of key-value pair in array by key
func (ol *OrderedList) getIndexByKey(key string) (int, error) {
	for i, record := range ol.Rec {
		if record.Key == key {
			return i, nil
		}
	}

	return -1, errors.New("key does not exists")
}
