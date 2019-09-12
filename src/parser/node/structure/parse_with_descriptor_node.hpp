#pragma once

#include <parser/node/parse_node.hpp>

namespace ql::parser {
    class ParseWithDescriptorNode : public ParseNode {
    protected:
        std::string_view m_InnerBody;
    public:
        ParseWithDescriptorNode(std::string&& body, std::string_view const& innerBody, std::vector<std::string>&& tokens, ParentRef const& parent)
                : ParseNode(std::move(body), std::move(tokens), parent), m_InnerBody(innerBody) {
        }
    };
}
