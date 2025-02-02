package elevatordriver


import (
	. "Project/dataenums"
	"encoding/json"
	"os"
)

func initelevator() (Elevator, bool) {
	elevator, err := loadElevator()
	if err == nil {
		return elevator, false // False means backup file was used
	}

	// Default initialization if file doesn't exist or decoding fails
	return Elevator{
		CurrentFloor:     -1, 
		Dirn:             MDDown,
		CurrentBehaviour: EBIdle,
		Config: ElevatorConfig{
			ClearRequestVariant: ClearRequestVariantConfig,
			DoorOpenDurationS:   DoorOpenDurationSConfig,
		},
	}, true // True means it was reset (no backup existed)
}



// SaveElevator saves the elevator state to "elevatorBackup.json"
func saveElevator(e Elevator) error {
	file, err := os.Create("elevatorBackup.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(e)
}

// LoadElevator loads the elevator state from "elevatorBackup.json"
func loadElevator() (Elevator, error) {
	file, err := os.Open("elevatorBackup.json")
	if err != nil {
		return Elevator{}, err
	}
	defer file.Close()

	var e Elevator
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&e); err != nil {
		return Elevator{}, err
	}
	return e, nil
}
