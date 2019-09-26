#pragma once

#include <parser/node/implementation/line_node.hpp>

namespace ql::parser {
    class VariableNode : public LineNode {
    public:
        using LineNode::LineNode;

        virtual uint getSize() {
            return 0u;
        }
    };
}
