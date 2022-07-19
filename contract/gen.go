package contract

//go:generate abigen --abi ./abi/ERC20.json --type ERC20 --pkg contract --out ./erc20.go
//go:generate abigen --abi ./abi/ERC721.json --type ERC721 --pkg contract --out ./erc721.go
//go:generate abigen --abi ./abi/ERC1155.json --type ERC1155 --pkg contract --out ./erc1155.go
