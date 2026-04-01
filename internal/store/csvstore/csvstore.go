package csvstore

import (
	"encoding/csv"
	"net"
	"os"
	"strings"

	"zivusomer/torq-project/internal/location"
	"zivusomer/torq-project/internal/store"
)

type Store struct {
	records map[string]location.Record
}

func New(path string) (*Store, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	records := make(map[string]location.Record, len(rows))
	for i, row := range rows {
		if len(row) < 3 {
			continue
		}

		ip := strings.TrimSpace(row[0])
		city := strings.TrimSpace(row[1])
		country := strings.TrimSpace(row[2])
		if ip == "" || city == "" || country == "" {
			continue
		}
		if i == 0 && strings.EqualFold(ip, "ip") && strings.EqualFold(city, "city") && strings.EqualFold(country, "country") {
			continue
		}

		records[ip] = location.Record{
			Country: country,
			City:    city,
		}
	}

	return &Store{records: records}, nil
}

func (s *Store) FindByIP(ip net.IP) (location.Record, error) {
	record, ok := s.records[ip.String()]
	if !ok {
		return location.Record{}, store.ErrIPNotFound
	}
	return record, nil
}
