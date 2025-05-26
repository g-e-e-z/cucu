cmake_minimum_required(VERSION 3.14)
project(cucu)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)


# JSON (release archive)
FetchContent_Declare(
  json
  URL https://github.com/nlohmann/json/releases/download/v3.12.0/json.tar.xz
)
FetchContent_MakeAvailable(json)

# FTXUI
FetchContent_Declare(
  ftxui
  GIT_REPOSITORY https://github.com/ArthurSonzogni/FTXUI
  GIT_TAG        main
)
FetchContent_MakeAvailable(ftxui)

add_executable(cucu
  main.cpp
  components/context_bar.cpp
)

target_include_directories(cucu PRIVATE components)

# Link FTXUI components
target_link_libraries(cucu PRIVATE ftxui::screen ftxui::dom ftxui::component)

# Link JSON correctly
target_link_libraries(cucu PRIVATE nlohmann_json)

