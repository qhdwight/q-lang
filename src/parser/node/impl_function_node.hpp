#pragma once

#include "parse_node.hpp"

namespace ql::parser {
    class ImplementFunctionNode : public ParseNode {
    public:
        using ParseNode::ParseNode;
    };
}
