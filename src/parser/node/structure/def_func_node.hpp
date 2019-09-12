#pragma once

#include <parser/node/structure/parse_with_descriptor_node.hpp>

namespace ql::parser {
    class DefineFunctionNode : public ParseWithDescriptorNode {
    public:
        using ParseWithDescriptorNode::ParseWithDescriptorNode;
    };
}
