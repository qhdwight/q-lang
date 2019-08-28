#include "parse_node.hpp"

namespace ql::parse {

    void ParseNode::addChild(const std::shared_ptr<ParseNode>& node) {
        m_Children.push_back(node);
    }
}