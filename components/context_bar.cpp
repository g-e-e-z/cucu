#include "context_bar.hpp"

#include <ftxui/component/component.hpp>
#include <ftxui/component/event.hpp>
#include <ftxui/dom/elements.hpp>
#include <string>
#include <vector>

using namespace ftxui;

Component MakeContextBar(std::string* url, int* method_index, std::vector<std::string>* methods) {
    auto method_dropdown = Dropdown(methods, method_index);

    auto url_input = Input(url, "Request URL") | flex | size(HEIGHT, EQUAL, 1);
    url_input |= CatchEvent([](Event event) {
        return event == Event::Return;
    });

    Component context_bar = Container::Horizontal({
        method_dropdown,
        url_input,
    });

    return Renderer(context_bar, [=] {
        return window(
            text("Context"),
            hbox({
                text(" ["), text((*methods)[*method_index]), text(" ▼] "),
                separator(), text(" "), text("URL: "), url_input->Render()
            })
        );
    });
}

