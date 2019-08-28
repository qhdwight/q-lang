#include "abstract_node.hpp"

namespace ql::parser {
    void AbstractNode::addChild(std::shared_ptr<AbstractNode> const& node) {
        m_Children.push_back(node);
    }
}
