#include "read.hpp"

#include <fstream>
#include <optional>

namespace ql::util {
    std::optional<std::string> readAllText(std::string const& fileName) {
        std::ifstream t(fileName);
        if (!t) return {};
        std::string text;
        t.seekg(0, std::ios::end);
        text.reserve(t.tellg());
        t.seekg(0, std::ios::beg);
        text.assign(std::istreambuf_iterator<char>(t), std::istreambuf_iterator<char>());
        return text;
    }
}
