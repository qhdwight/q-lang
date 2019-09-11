#pragma once

#include <parser/node/parse_node.hpp>

namespace ql::parser {
    class LineNode : public ParseNode {
    public:
        using ParseNode::ParseNode;
    };
}
