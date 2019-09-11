#include "package_node.hpp"

namespace ql::parser {
    void PackageNode::parse() {
        m_Name = m_Tokens[1];
    }
}
