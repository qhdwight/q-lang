#include "read.hpp"

#include <fstream>

namespace ql::util {
    std::optional<std::string> readAllText(std::string const& fileName) {
        std::ifstream fileStream(fileName);
        if (!fileStream) return {};
        std::string text;
        fileStream.seekg(0, std::ios::end);
        text.reserve(fileStream.tellg());
        fileStream.seekg(0, std::ios::beg);
        text.assign(std::istreambuf_iterator<char>(fileStream), std::istreambuf_iterator<char>());
        return text;
    }
}
