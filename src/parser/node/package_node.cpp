#include "package_node.hpp"

namespace ql::parser {
    void PackageNode::parse(std::string const& text, std::vector<std::string> const& tokens) {
        m_Name = tokens[1];
    }
}
