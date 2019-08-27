#pragma once

#include <string>
#include <optional>

namespace ql::util {
    std::optional<std::string> readAllText(std::string const& fileName);
}
