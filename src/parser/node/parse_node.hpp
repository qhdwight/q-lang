#pragma once

#include <parser/node/abstract_node.hpp>

namespace ql::parser {
    typedef std::vector<std::string> Tokens;

    class ParseNode : public AbstractNode {
    protected:
        std::string m_Body;
        Tokens m_Tokens;
    public:
        ParseNode(std::string&& body, Tokens&& tokens, ParentRef const& parent)
                : AbstractNode(parent), m_Body(body), m_Tokens(tokens) {
        };

        std::string_view getText() const { return m_Body; }

        virtual void parse() {}
    };
}
