#pragma once

#include <parser/node/implementation/variable_node.hpp>

namespace ql::parser {
    class IntegerNode : public VariableNode {
    public:
        using VariableNode::VariableNode;

        void parse() override;

        uint getSize() override;
    };
}
