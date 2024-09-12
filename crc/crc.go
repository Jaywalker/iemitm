package crc

type CRC struct {
	m_dwCRC32 [256]uint64
}

func (this *CRC) Calculate(data []byte, size uint64) uint64 {
	var dwSize uint64
	dwSize = size
	pData := data
	var dwCrc uint64
	dwCrc = 0
	var k uint64
	k = 0

	// This is where the CRC lives. Null it out before we calculate it
	pData[6] = 0x00
	pData[7] = 0x00
	pData[8] = 0x00
	pData[9] = 0x00

	for k = 0; k < dwSize; k++ {
		dataByteAsUint := uint64(pData[k])
		dwCrc = this.m_dwCRC32[(dwCrc&0xFF)^dataByteAsUint] ^ (dwCrc >> 8)
	}
	return dwCrc
}

func New() (ret *CRC) {
	var m_dwCRC32 [256]uint64
	// Init our CRC values
	var m uint64
	m = 0
	for m = 0; m <= 255; m++ {
		v1 := m
		var bit uint64
		bit = 0
		for bit = 8; bit > 0; bit-- {
			if v1&1 != 0 {
				v1 >>= 1
				v1 ^= 0xEDB88320
			} else {
				v1 >>= 1
			}
		}
		m_dwCRC32[m] = v1
	}
	ret = &CRC{m_dwCRC32}
	return
}
