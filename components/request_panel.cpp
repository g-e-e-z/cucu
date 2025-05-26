#include "request_panel.hpp"

#include <ftxui/component/component.hpp>
#include <ftxui/component/event.hpp>
#include <ftxui/dom/elements.hpp>
#include <string>
#include <vector>

using namespace ftxui;

Component MakeRequestPanel(
    std::vector<std::pair<std::string, std::string>>* headers) {
    std::vector<Element> header_elements;
    for (auto& [key, value] : *headers) {
        header_elements.push_back(
            hbox({text(key) | bold | flex, text(": "), text(value)}));
    };

    return Renderer(
        [=] { return window(text("Request"), vbox(header_elements) | flex); });
}
