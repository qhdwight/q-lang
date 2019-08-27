#pragma once

#include <memory>
#include <vector>
#include <utility>

namespace ql::parse {
    class ParseNode {
        typedef std::vector<std::shared_ptr<ParseNode>> ParseNodeChildren;
        typedef std::weak_ptr<ParseNode> ParseNodeParent;

    private:
        ParseNodeChildren m_Children;
        ParseNodeParent m_Parent;
    public:
        ParseNode(ParseNodeParent parent, ParseNodeChildren children) : m_Parent(std::move(parent)), m_Children(std::move(children)) {}
    };
}
