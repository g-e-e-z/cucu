#include <ftxui/component/component.hpp>
#include <ftxui/component/component_base.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <ftxui/dom/elements.hpp>
#include <iostream>
#include <nlohmann/json.hpp>
#include <string>
#include <vector>
using namespace ftxui;

using json = nlohmann::json;

int main() {
    auto screen = ScreenInteractive::Fullscreen();

    // --- Context Bar ---
    std::vector<std::string> methods = {"GET", "POST", "PUT", "DELETE"};
    int selected_method = 0;
    auto method_dropdown = Dropdown(&methods, &selected_method);

    std::string url;  //= "https://api.example.com";
    auto url_input = Input(&url, "Request URL") | flex | size(HEIGHT, EQUAL, 1);
    url_input |= CatchEvent([&](Event event) { return event == Event::Return; });

    Component context_bar = Container::Horizontal({
        method_dropdown,
        url_input,
    });

    auto context_renderer = Renderer(context_bar, [&] {
        return window(
            text("Context"),
            hbox({text(" ["), text(methods[selected_method]), text(" ▼] "),
                  separator(), text(" "), text("URL: "), url_input->Render()}));
    });

    // --- Headers (Placeholder) ---
    std::vector<std::pair<std::string, std::string>> headers = {
        {"Authorization", "Bearer <token>"},
        {"Accept", "application/json"},
    };

    std::vector<Element> header_elements;
    for (auto& [key, value] : headers) {
        header_elements.push_back(
            hbox({text(key) | bold | flex, text(": "), text(value)}));
        // hbox({text(key) | bold | flex, text(": "), text(value) | flex}));
    }

    auto headers_renderer = Renderer(
        [&] { return window(text("Request"), vbox(header_elements) | flex); });

    // --- Response (Placeholder) ---
    std::string raw_json =
        R"({"user":{"id":42,"name":"alex","email":"a@b.com"}})";

    std::cout << "Raw JSON:\n" << raw_json << "\n";

    std::string formatted_json;
    try {
        auto parsed = json::parse(raw_json);
        formatted_json = parsed.dump(2);  // 2-space indentation
        std::cout << "Formatted JSON:\n" << formatted_json << "\n";
    } catch (const std::exception& e) {
        formatted_json = std::string("Failed to parse JSON: ") + e.what();
        std::cerr << "JSON parse error: " << e.what() << "\n";
    }
    auto response_renderer = Renderer([&] {
        return window(
            text("Response"),
            vbox({text("200 OK | 82ms | 2.4KB | application/json") | bold,
                  separator(), paragraph(formatted_json) | flex}) |
                flex);
    });

    // Toggle between request/response views
    bool show_response = false;  // Flip this manually for now
    auto center_pane = Renderer([&] {
        return show_response ? response_renderer->Render()
                             : headers_renderer->Render();
    });

    // --- Footer ---
    auto footer = Renderer([] {
        return text(
                   "[i] Insert  [x] Delete  [j/k] Navigate  [:] Command  [/] "
                   "Filter") |
               border;
    });

    // --- Combine Layout ---
    Component root = Container::Vertical({context_bar, center_pane, footer});

    auto app = Renderer(root, [&] {
        return vbox({context_renderer->Render(), center_pane->Render() | flex,
                     footer->Render()});
    });
    screen.Loop(app);
    // auto layout = Renderer([&] {
    //     return vbox({context_renderer->Render(), center_pane->Render() |
    //     flex,
    //                  footer->Render()});
    // });
    // screen.Loop(layout);
}
