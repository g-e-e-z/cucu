#pragma once
#include <ftxui/component/component.hpp>
#include <string>
#include <vector>

ftxui::Component MakeContextBar(std::string* url, int* method_index, std::vector<std::string>* methods);

