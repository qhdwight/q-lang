#pragma once

#include "parse_node.hpp"

namespace ql::parser {
    class PackageNode : public ParseNode {
    private:
        std::string m_Name;
    public:
        PackageNode(std::string&& rawText, ParentRef const& parent)
                : ParseNode(std::move(rawText), parent) {}

        void parse(std::string const& text) override;
    };
}