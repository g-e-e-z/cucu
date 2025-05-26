# cucu

# Build Instructions

## Initialization

```bash
mkdir build
cd build
cmake ..
make
./cucu

```

## Development

1. Make code changes

2. Rebuild the project
```bash
cd build
make
```
3. Run the binary: `./cucu`


## Clean & Build

```bash
rm -rf build
mkdir build
cd build
cmake ..
make

```

# CMakeLists Notes
After adding in the json package, clang not recognizing the package. Fixed with the following:
`cmake -S . -B build -DCMAKE_EXPORT_COMPILE_COMMANDS=ON`

