#pragma once

#include <parser/node/structure/parse_with_descriptor_node.hpp>

namespace ql::parser {
    class ImplementFunctionNode : public ParseWithDescriptorNode {
    public:
        using ParseWithDescriptorNode::ParseWithDescriptorNode;

        void parse() override;
    };
}
