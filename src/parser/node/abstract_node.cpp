#include "abstract_node.hpp"

namespace ql::parser {
    void AbstractNode::addChild(const std::shared_ptr<AbstractNode>& node) {
        m_Children.push_back(node);
    }
}
