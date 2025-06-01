#pragma once
#include <ftxui/component/component.hpp>
#include <string>
#include <vector>

ftxui::Component MakeRequestPanel(std::vector<std::pair<std::string, std::string>>* headers);

