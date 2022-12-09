package lib

import "time"

type GetSensorReadingsRawResponse struct {
	Data []interface{} `json:"data"`
}

type GetSensorLastReadingRawResponse struct {
	Data interface{} `json:"data"`
}

type SensorReading struct {
	Scap                   float64   `json:"scap"`
	Vcap                   float64   `json:"vcap"`
	Schema                 string    `json:"schema"`
	Timestamp              int64     `json:"timestamp"`
	MacAddress             string    `json:"macAddress"`
	RssiUplink             int       `json:"rssiUplink"`
	ReadingDate            time.Time `json:"readingDate"`
	Unsolicited            bool      `json:"unsolicited"`
	SchemaVersion          string    `json:"schemaVersion"`
	GatewaySerialNumber    string    `json:"gatewaySerialNumber"`
	PacketNumberGateway    int       `json:"packetNumberGateway"`
	PacketNumberEversensor int       `json:"packetNumberEversensor"`
}

type GetEversensorsShortResponse struct {
	Data []struct {
		Customer struct {
			ID        string `json:"id"`
			IsSandbox bool   `json:"isSandbox"`
			Name      string `json:"name"`
			Status    string `json:"status"`
		} `json:"customer"`
		DevkitBundled   bool `json:"devkitBundled"`
		LastAssociation struct {
			GatewaySerialNumber string    `json:"gatewaySerialNumber"`
			Timestamp           time.Time `json:"timestamp"`
		} `json:"lastAssociation"`
		LastInfo struct {
			FirmwareVersion     string    `json:"firmwareVersion"`
			GatewaySerialNumber string    `json:"gatewaySerialNumber"`
			PartNumber          string    `json:"partNumber"`
			SensorSerialNumber  string    `json:"sensorSerialNumber"`
			Timestamp           time.Time `json:"timestamp"`
		} `json:"lastInfo"`
		MacAddress                  string `json:"macAddress"`
		ManufacturedFirmwareVersion string `json:"manufacturedFirmwareVersion"`
		Type                        string `json:"type"`
	} `json:"data"`
	PaginationInfo struct {
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
		TotalItems int `json:"totalItems"`
		TotalPages int `json:"totalPages"`
	} `json:"paginationInfo"`
}
