package utils

import "github.com/gofrs/uuid"

type PartitionUUID struct {
	Persistent uuid.UUID
	Oem        uuid.UUID
	State      uuid.UUID
	Recovery   uuid.UUID
	Bios       uuid.UUID
	Efi        uuid.UUID
}

func ReturnPartitionFixedUUIDs() *PartitionUUID {
	return &PartitionUUID{
		Persistent: uuid.NewV5(uuid.NamespaceURL, "KAIROS_PERSISTENT"),
		Oem:        uuid.NewV5(uuid.NamespaceURL, "KAIROS_OEM"),
		State:      uuid.NewV5(uuid.NamespaceURL, "KAIROS_STATE"),
		Recovery:   uuid.NewV5(uuid.NamespaceURL, "KAIROS_RECOVERY"),
		Bios:       uuid.NewV5(uuid.NamespaceURL, "KAIROS_BIOS"),
		Efi:        uuid.NewV5(uuid.NamespaceURL, "KAIROS_EFI"),
	}
}
