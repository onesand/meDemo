package contract

import "github.com/ethereum/go-ethereum/common"

var (
	ZeroAddress        = common.HexToAddress("0x0000000000000000000000000000000000000000")
	NativePaymentToken = common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
)

var (
	InterfaceIDERC721              = [4]byte{0x80, 0xac, 0x58, 0xcd}
	InterfaceIDERC721TokenReceiver = [4]byte{0x15, 0x0b, 0x7a, 0x02}
	InterfaceIDERC721Metadata      = [4]byte{0x5b, 0x5e, 0x13, 0x9f}
	InterfaceIDERC721Enumerable    = [4]byte{0x78, 0x0e, 0x9d, 0x63}
)

var (
	InterfaceIDERC1155              = [4]byte{0xd9, 0xb6, 0x7a, 0x26}
	InterfaceIDERC1155TokenReceiver = [4]byte{0x4e, 0x23, 0x12, 0xe0}
	InterfaceIDERC1155Metadata_URI  = [4]byte{0x0e, 0x89, 0x34, 0x1c}
)

var (
	EnsContractAddress = common.HexToAddress("0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85")
)
