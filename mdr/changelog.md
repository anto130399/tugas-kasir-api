# Changelog

## [Unreleased]

### Refactor
- Renaming `models/produk_service.go` to `models/produk.go` to match content.
- Renaming `handlers/produk_hadler.go` to `handlers/produk_handler.go` to fix typo.
- Splitting `Category` logic from `Produk` files into dedicated files.

### New Features
- **Category API**: Implemented full CRUD for categories.
    - `repositories/category_repository.go`
    - `services/category_service.go`
    - `handlers/category_handler.go`
    - New routes in `main.go`.

### Errors Found & Fixed
- **Naming**: `models/produk_service.go` was misnamed (it contained structs).
- **Typo**: `handlers/produk_hadler.go` (missing 'n').
- **Mixed Concerns**: Category logic was tightly coupled with Product logic.
