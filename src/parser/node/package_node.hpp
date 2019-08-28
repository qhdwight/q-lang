#pragma once

#include "parse_node.hpp"

namespace ql::parser {
    class PackageNode : public ParseNode {
    private:
        std::string m_Name;
    public:
        using ParseNode::ParseNode;

        void parse(std::string const& text, std::vector<std::string> const& tokens) override;
    };
}