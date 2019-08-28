#pragma once

#include "parse_node.hpp"

namespace ql::parser {
    class DefineFunctionNode : public ParseNode {
    public:
        using ParseNode::ParseNode;

        void parse(std::string const& text, std::vector<std::string> const& tokens) override;
    };
}
